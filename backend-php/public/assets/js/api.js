(() => {
  // frontend/core/parse-data-attr.mts
  var assertString = (value, key) => {
    if (typeof value !== "string") {
      throw new Error(`data-fb-api.${key} must be a string.`);
    }
  };
  var assertBoolean = (value, key) => {
    if (typeof value !== "boolean") {
      throw new Error(`data-fb-api.${key} must be a boolean.`);
    }
  };
  var assertPositiveNumber = (value, key) => {
    if (typeof value !== "number" || !Number.isFinite(value) || value <= 0) {
      throw new Error(`data-fb-api.${key} must be a positive number.`);
    }
  };
  var assertPlainObject = (value, key) => {
    if (typeof value !== "object" || value === null || Array.isArray(value)) {
      throw new Error(`data-fb-api.${key} must be an object.`);
    }
  };
  var TOP_LEVEL_SCHEMA = {
    href: assertString,
    type: assertString,
    endpoint: assertString,
    callback: assertString,
    message: assertString,
    redirect: assertString,
    reload: assertBoolean,
    preventNavigation: assertBoolean,
    timeoutMs: assertPositiveNumber,
    timeoutMessage: assertString,
    params: assertPlainObject
  };
  var LOADING_STRING_FIELDS = ["message", "button", "target", "alertClass"];
  var MODAL_ALLOWED_TYPES = ["confirm", "danger", "prompt"];
  var MODAL_STRING_FIELDS = ["title", "content", "button", "buttonColor", "label", "value", "key"];
  function validateLoading(loading) {
    assertPlainObject(loading, "loading");
    for (const field of LOADING_STRING_FIELDS) {
      if (Object.prototype.hasOwnProperty.call(loading, field)) {
        assertString(loading[field], `loading.${field}`);
      }
    }
  }
  function validateModal(modal) {
    assertPlainObject(modal, "modal");
    if (typeof modal.type !== "string") {
      throw new Error("data-fb-api.modal.type must be a string.");
    }
    if (!MODAL_ALLOWED_TYPES.includes(modal.type)) {
      throw new Error(`data-fb-api.modal.type must be one of: ${MODAL_ALLOWED_TYPES.join(", ")}.`);
    }
    if (modal.type === "prompt" && typeof modal.key !== "string") {
      throw new Error("data-fb-api.modal.key is required for prompt modals.");
    }
    for (const field of MODAL_STRING_FIELDS) {
      if (Object.prototype.hasOwnProperty.call(modal, field)) {
        assertString(modal[field], `modal.${field}`);
      }
    }
  }
  function parseDataAttr(dataAttrValue) {
    if (!dataAttrValue) {
      return {};
    }
    let data;
    try {
      data = JSON.parse(dataAttrValue);
    } catch (error) {
      throw new Error("Invalid JSON in data-fb-api attribute.");
    }
    if (typeof data !== "object" || data === null || Array.isArray(data)) {
      throw new Error("data-fb-api must be a JSON object.");
    }
    for (const [key, validator] of Object.entries(TOP_LEVEL_SCHEMA)) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        validator(data[key], key);
      }
    }
    if (Object.prototype.hasOwnProperty.call(data, "loading")) {
      validateLoading(data.loading);
    }
    if (Object.prototype.hasOwnProperty.call(data, "modal")) {
      validateModal(data.modal);
    }
    return data;
  }

  // frontend/core/api-helpers.mts
  function isApiResponsePayload(payload) {
    return payload === null || typeof payload === "object" && !Array.isArray(payload);
  }
  function injectCSRFToken(params, token) {
    if (params instanceof FormData) {
      if (!params.has("CSRFToken")) {
        params.append("CSRFToken", token);
      }
    } else if (params && typeof params === "object") {
      if (!params.CSRFToken) {
        params.CSRFToken = token;
      }
    }
    return params;
  }
  function buildRequestBody(method, params, url) {
    const methodLower = method.toLowerCase();
    const isFormData = params instanceof FormData;
    let body = null;
    if (methodLower === "get") {
      if (isFormData) {
        for (const [key, value] of params.entries()) {
          if (key !== "CSRFToken") {
            url.searchParams.append(key, String(value));
          }
        }
      } else if (params && typeof params === "object") {
        Object.keys(params).filter((key) => key !== "CSRFToken").forEach((key) => url.searchParams.append(key, String(params[key])));
      } else if (params) {
        url.search = String(params);
        url.searchParams.delete("CSRFToken");
      }
    } else if (["post", "put", "patch", "delete"].includes(methodLower)) {
      if (isFormData || typeof params === "string") {
        body = params;
      } else {
        body = JSON.stringify(params);
      }
    }
    return { body, isFormData };
  }
  function buildHeaders({ url, body, isFormData, csrfToken, origin }) {
    const headers = {
      "Accept": "application/json"
    };
    if (url.origin === origin) {
      headers["X-Requested-With"] = "XMLHttpRequest";
      headers["X-CSRF-Token"] = csrfToken || "";
    }
    if (body && !isFormData) {
      headers["Content-Type"] = "application/json";
    }
    return headers;
  }
  async function parseResponseBody(response) {
    var _a, _b;
    if (response.status === 204) {
      return { payload: null, rawText: "" };
    }
    const text = (_b = await ((_a = response.text) == null ? void 0 : _a.call(response))) != null ? _b : "";
    if (!text) {
      return { payload: null, rawText: "" };
    }
    try {
      return { payload: JSON.parse(text), rawText: text };
    } catch (e) {
      return { payload: null, rawText: text };
    }
  }
  function validateHttpResponse(response, parsed) {
    var _a, _b;
    if (!response.ok) {
      const payload = isApiResponsePayload(parsed.payload) ? parsed.payload : null;
      const error = new Error(((_a = payload == null ? void 0 : payload.error) == null ? void 0 : _a.message) || `HTTP error ${response.status}: ${response.statusText}`);
      error.code = ((_b = payload == null ? void 0 : payload.error) == null ? void 0 : _b.code) || `http_${response.status}`;
      error.status = response.status;
      error.rawBody = parsed.rawText;
      throw error;
    }
    if (parsed.rawText && parsed.payload === null) {
      throw new Error("Invalid or non-JSON response from server");
    }
    return isApiResponsePayload(parsed.payload) ? parsed.payload : null;
  }
  function interpretResponse(payload) {
    if (!payload) {
      return null;
    }
    if (payload.error) {
      const error = new Error(payload.error.message || "Unknown API error");
      error.code = payload.error.code;
      throw error;
    }
    return payload.result;
  }
  function normalizeApiError(error, { timeoutMs = 3e4, timeoutMessage = null } = {}) {
    if (error.name === "AbortError") {
      return {
        message: timeoutMessage || `Request timed out after ${timeoutMs / 1e3} seconds`,
        code: "timeout_error"
      };
    }
    if (error.name === "TypeError" && /NetworkError|Failed to fetch|Load failed/i.test(error.message)) {
      return {
        message: "Network connection error",
        code: "network_error"
      };
    }
    return {
      message: error.message || "Unknown error occurred",
      code: error.code || "unknown_error"
    };
  }

  // frontend/core/link-helpers.mts
  function dispatchLinkAction(apiData, rawHref, modalsLib, onRequest) {
    var _a, _b, _c, _d, _e, _f, _g, _h;
    if (!apiData.hasOwnProperty("modal")) {
      onRequest("GET", rawHref);
      return;
    }
    const modal = apiData.modal;
    if (!modalsLib || typeof modalsLib.create !== "function") {
      if (modal.type === "prompt") {
        const value = window.prompt((_b = (_a = modal.label) != null ? _a : modal.title) != null ? _b : "", (_c = modal.value) != null ? _c : "");
        if (value) {
          onRequest("GET", rawHref, { [modal.key]: value });
        }
      } else if (window.confirm(modal.content || modal.title || "Are you sure?")) {
        onRequest("GET", rawHref);
      }
    } else if (modal.type === "prompt") {
      modalsLib.create({
        type: modal.type,
        title: modal.title,
        label: (_d = modal.label) != null ? _d : "Label",
        value: (_e = modal.value) != null ? _e : "",
        promptConfirmCallback: (value) => {
          if (value) {
            onRequest("GET", rawHref, { [modal.key]: value });
          }
        }
      });
    } else {
      modalsLib.create({
        type: modal.type === "confirm" ? "small-confirm" : modal.type,
        title: modal.title,
        content: (_f = modal.content) != null ? _f : "",
        confirmButton: (_g = modal.button) != null ? _g : "Confirm",
        confirmButtonColor: (_h = modal.buttonColor) != null ? _h : "primary",
        confirmCallback: () => {
          onRequest("GET", rawHref);
        }
      });
    }
  }
  function createLinkLoadingState(linkElement) {
    let requestInProgress = false;
    let loadingAlert = null;
    let beforeUnloadHandler = null;
    let originalHtml = null;
    let originalAriaBusy = null;
    let originalAriaDisabled = null;
    let originallyDisabled = false;
    const getLoadingTarget = (selector) => {
      if (selector) {
        try {
          const target = document.querySelector(selector);
          if (target) {
            return target;
          }
        } catch (error) {
          console.warn("Invalid loading target selector:", selector);
        }
      }
      return linkElement.closest(".card-footer") || linkElement.parentElement;
    };
    const set = (apiData) => {
      var _a, _b;
      requestInProgress = true;
      originalHtml = linkElement.innerHTML;
      originalAriaBusy = linkElement.getAttribute("aria-busy");
      originalAriaDisabled = linkElement.getAttribute("aria-disabled");
      originallyDisabled = linkElement.classList.contains("disabled");
      linkElement.setAttribute("aria-busy", "true");
      linkElement.setAttribute("aria-disabled", "true");
      linkElement.classList.add("disabled");
      if ((_a = apiData.loading) == null ? void 0 : _a.button) {
        const spinner = document.createElement("span");
        spinner.className = "spinner-border spinner-border-sm me-2";
        spinner.setAttribute("aria-hidden", "true");
        linkElement.replaceChildren(spinner, document.createTextNode(apiData.loading.button));
      }
      if ((_b = apiData.loading) == null ? void 0 : _b.message) {
        const target = getLoadingTarget(apiData.loading.target);
        if (target) {
          loadingAlert = document.createElement("div");
          loadingAlert.className = apiData.loading.alertClass || "alert alert-info mt-3 mb-0";
          loadingAlert.setAttribute("role", "status");
          loadingAlert.textContent = apiData.loading.message;
          target.appendChild(loadingAlert);
        }
      }
      if (apiData.preventNavigation) {
        beforeUnloadHandler = (event) => {
          event.preventDefault();
          event.returnValue = "";
        };
        window.addEventListener("beforeunload", beforeUnloadHandler);
      }
    };
    const reset = () => {
      if (!requestInProgress) {
        return;
      }
      requestInProgress = false;
      if (originalHtml !== null) {
        linkElement.innerHTML = originalHtml;
      }
      if (originalAriaBusy === null) {
        linkElement.removeAttribute("aria-busy");
      } else {
        linkElement.setAttribute("aria-busy", originalAriaBusy);
      }
      if (originalAriaDisabled === null) {
        linkElement.removeAttribute("aria-disabled");
      } else {
        linkElement.setAttribute("aria-disabled", originalAriaDisabled);
      }
      if (!originallyDisabled) {
        linkElement.classList.remove("disabled");
      }
      if (loadingAlert) {
        loadingAlert.remove();
        loadingAlert = null;
      }
      if (beforeUnloadHandler) {
        window.removeEventListener("beforeunload", beforeUnloadHandler);
        beforeUnloadHandler = null;
      }
    };
    return {
      set,
      reset,
      isInProgress: () => requestInProgress
    };
  }

  // frontend/core/api.ts
  var FOSSBilling = window.FOSSBilling = window.FOSSBilling || {};
  var Tools = {
    /**
     * Constructs the full URL for an API endpoint.
     * If the provided URL is relative, it's resolved against the application's base API URL.
     *
     * @param {string} url The API endpoint (e.g., "guest/system/company") or a full URL.
     * @returns {string} The complete URL for the API call.
     */
    getBaseURL: function(url) {
      if (typeof url !== "string" || !url.trim()) {
        return `${window.location.origin}/api/`;
      }
      if (url.startsWith("http://") || url.startsWith("https://")) {
        return url;
      }
      if (url.includes("index.php?_url=/api/") || url.includes("?_url=/api/")) {
        return new URL(url, `${window.location.origin}/`).toString();
      }
      const base = `${window.location.origin}/api/`;
      let normalized = url;
      if (normalized.startsWith("/")) {
        normalized = normalized.slice(1);
      }
      if (normalized.startsWith("api/")) {
        normalized = normalized.slice(4);
      }
      return new URL(normalized, base).toString();
    },
    /**
     * @returns {string|null} The CSRF token from cookie, or null if not found.
     */
    getCSRFToken: function() {
      const match = document.cookie.match(/(?:^|;\s*)fossbilling_csrf=([^;]*)/) || document.cookie.match(/(?:^|;\s*)csrf_token=([^;]*)/);
      return match ? decodeURIComponent(match[1]) : null;
    },
    /**
     * Check if a string is valid JSON or not.
     *
     * @param {string} jsonString The string to check.
     * @returns {boolean} True if the string is valid JSON, or false if it is not.
     */
    isJSON: function(jsonString) {
      if (typeof jsonString !== "string") {
        return false;
      }
      try {
        JSON.parse(jsonString);
        return true;
      } catch (error) {
        return false;
      }
    },
    /**
     * Converts a FormData object into a urlencoded string.
     *
     * @param {FormData} formData The FormData object to serialize.
     * @returns {string} Serialized string of the FormData.
     */
    serializeFormData: function(formData) {
      const params = new URLSearchParams(formData);
      if (!formData.has("CSRFToken")) {
        const token = Tools.getCSRFToken();
        if (token) {
          params.append("CSRFToken", token);
        }
      }
      return params.toString();
    },
    /**
     * Converts a FormData object into a valid object.
     *
     * @param {FormData} formData The FormData object to serialize.
     * @returns {object} The reformatted object.
     */
    serializeFormDataToObject: function(formData) {
      const obj = {};
      for (const [key, value] of formData.entries()) {
        if (key.endsWith("[]")) {
          const plainKey = key.slice(0, -2);
          if (!obj[plainKey]) {
            obj[plainKey] = [];
          }
          obj[plainKey].push(value);
        } else if (Object.prototype.hasOwnProperty.call(obj, key)) {
          obj[key] = value;
        } else {
          obj[key] = value;
        }
      }
      if (!Object.prototype.hasOwnProperty.call(obj, "CSRFToken")) {
        const token = Tools.getCSRFToken();
        if (token) {
          obj.CSRFToken = token;
        }
      }
      const reformattedObj = {};
      Object.keys(obj).forEach((originalKey) => {
        const parts = originalKey.match(/[^[\]]+/g) || [originalKey];
        let currentContext = reformattedObj;
        for (let i = 0; i < parts.length; i++) {
          const part = parts[i];
          if (i === parts.length - 1) {
            currentContext[part] = obj[originalKey];
          } else {
            if (!Object.prototype.hasOwnProperty.call(currentContext, part) || typeof currentContext[part] !== "object" || currentContext[part] === null) {
              currentContext[part] = {};
            }
            currentContext = currentContext[part];
          }
        }
      });
      return reformattedObj;
    },
    /**
     * Converts a FormData object into a valid JSON string.
     *
     * @param {FormData} formData The FormData object to serialize.
     * @returns {string} JSON string of the FormData object.
     */
    serializeFormDataToJSON: function(formData) {
      return JSON.stringify(this.serializeFormDataToObject(formData));
    }
  };
  Tools.parseDataAttr = parseDataAttr;
  function _createApiRole(role) {
    const baseNamespaceUrlString = Tools.getBaseURL(role);
    const createMethod = (method) => {
      return function(endpoint, params, successHandler, errorHandler, enableLoader = true) {
        if (typeof endpoint !== "string" || !endpoint.trim()) {
          throw new Error("Invalid endpoint: must be a non-empty string");
        }
        const requestUrl = new URL(endpoint, `${baseNamespaceUrlString}/`).toString();
        API.makeRequest(method, requestUrl, params, successHandler, errorHandler, enableLoader);
      };
    };
    return {
      baseURL: baseNamespaceUrlString,
      get: createMethod("GET"),
      post: createMethod("POST"),
      put: createMethod("PUT"),
      delete: createMethod("DELETE"),
      patch: createMethod("PATCH")
    };
  }
  var API = {
    /**
     * Wrapper for the admin API.
     * @documentation https://docs.fossbilling.org/extensions-and-development/javascript/
     */
    admin: _createApiRole("admin"),
    /**
     * Wrapper for the client API.
     * @documentation https://docs.fossbilling.org/extensions-and-development/javascript/
     */
    client: _createApiRole("client"),
    /**
     * Wrapper for the guest API.
     * @documentation https://docs.fossbilling.org/extensions-and-development/javascript/
     */
    guest: _createApiRole("guest"),
    /**
     * Make a request to the API.
     *
     * @param {string} method The HTTP method to use.
     * @param {string} url The URL to call.
     * @param {object|string} [params] The parameters to send.
     * @param {function} [successHandler] The function to call if the request is successful.
     * @param {function} [errorHandler] The function to call if the request is unsuccessful.
     * @param {boolean} [enableLoader=true] Enable or disable the usage of a loader. Custom themes simply need to provide one with the spinner-border class.
     * @param {number} [timeoutMs=30000] Timeout duration in milliseconds.
     * @param {string|null} [timeoutMessage=null] Message to show when the request times out.
     * @documentation https://docs.fossbilling.org/extensions-and-development/javascript/
     */
    makeRequest: function(method, url, params, successHandler, errorHandler, enableLoader = true, timeoutMs = 3e4, timeoutMessage = null) {
      let loader = enableLoader ? this._createLoader() : null;
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), timeoutMs);
      const csrfToken = Tools.getCSRFToken();
      const urlObj = new URL(url);
      const isSameOrigin = urlObj.origin === window.location.origin;
      if (isSameOrigin) {
        injectCSRFToken(params, csrfToken);
      }
      const { body, isFormData } = buildRequestBody(method, params, urlObj);
      const headers = buildHeaders({ url: urlObj, body, isFormData, csrfToken, origin: window.location.origin });
      const fetchOptions = {
        method,
        headers,
        signal: controller.signal
      };
      if (method.toLowerCase() !== "get") {
        fetchOptions.body = body;
      }
      return fetch(urlObj.toString(), fetchOptions).then(async (response) => {
        clearTimeout(timeoutId);
        if (response.redirected) {
          window.location.replace(response.url);
          return;
        }
        const parsed = await parseResponseBody(response);
        return validateHttpResponse(response, parsed);
      }).then((payload) => {
        const result = interpretResponse(payload);
        if (typeof successHandler === "function") {
          successHandler(result);
        }
        return result;
      }).catch((error) => {
        clearTimeout(timeoutId);
        const errorObj = normalizeApiError(error, { timeoutMs, timeoutMessage });
        console.error(`API Error: ${errorObj.message}`);
        if (typeof errorHandler === "function") {
          errorHandler(errorObj);
        } else {
          console.warn("No error handler was specified for API error.");
          const normalizedError = new Error(errorObj.message);
          normalizedError.code = errorObj.code;
          throw normalizedError;
        }
      }).finally(() => {
        if (enableLoader && loader) {
          if (loader._fadeInTimeout) {
            clearTimeout(loader._fadeInTimeout);
          }
          if (document.body.contains(loader)) {
            document.body.removeChild(loader);
          }
          loader = null;
        }
      });
    },
    /**
     * After the API request is complete, this function will be called.
     *
     * @param {object} object The HTML element that triggered the API call.
     * @param {*} result The result of the API call.
     * @returns
     */
    _afterComplete: function(object, result) {
      let apiData;
      try {
        apiData = Tools.parseDataAttr(object.dataset.fbApi || "{}");
      } catch (error) {
        console.warn("Invalid JSON in data-fb-api attribute:", error);
        return;
      }
      if (apiData.hasOwnProperty("callback") && typeof window[apiData.callback] === "function") {
        return window[apiData.callback](result);
      } else if (apiData.hasOwnProperty("callback")) {
        console.warn("Invalid callback function:", apiData.callback);
      }
      if (apiData.hasOwnProperty("redirect")) {
        const redirectUrl = new URL(apiData.redirect, window.location.href);
        if (redirectUrl.href === window.location.href) {
          window.location.reload();
        } else {
          window.location = redirectUrl.href;
        }
        return;
      }
      if (apiData.hasOwnProperty("reload")) {
        window.location.reload();
        return;
      }
      if (apiData.hasOwnProperty("message")) {
        FOSSBilling.ui.notify(apiData.message, "success");
        return;
      }
      if (result) {
        FOSSBilling.ui.notify("Form Updated", "success");
        return;
      }
      console.warn("Unhandled API response in _afterComplete:", apiData);
    },
    /**
     * Creates a loader element and appends it to the document body.
     *
     * @returns {HTMLElement} The created loader element.
     */
    _createLoader: function() {
      const loader = document.createElement("div");
      loader.classList.add("spinner-border");
      loader.setAttribute("role", "status");
      Object.assign(loader.style, {
        width: "4rem",
        height: "4rem",
        left: "50%",
        top: "50%",
        position: "fixed",
        opacity: "0",
        transition: "opacity 250ms"
      });
      document.body.appendChild(loader);
      loader._fadeInTimeout = setTimeout(() => {
        loader.style.opacity = "1";
      }, 250);
      return loader;
    },
    /**
     * Attach event listeners to forms with data attribute 'data-fb-api'.
     **/
    _apiForm: function() {
      const formElements = document.querySelectorAll("form[data-fb-api]");
      if (formElements.length > 0) {
        formElements.forEach((formElement) => {
          if (formElement.dataset.fbApiBound === "true") {
            return;
          }
          formElement.dataset.fbApiBound = "true";
          formElement.addEventListener("submit", function(event) {
            event.preventDefault();
            const formData = new FormData(formElement);
            const submitter = event.submitter;
            if (submitter == null ? void 0 : submitter.name) {
              formData.append(submitter.name, submitter.value);
            }
            if (FOSSBilling.editor) {
              if (!FOSSBilling.editor.validateForm(formElement)) {
                return FOSSBilling.ui.notify("At least one of the required fields are empty.", "error");
              }
              FOSSBilling.editor.syncForm(formElement, formData);
            }
            const formMethod = (formElement.getAttribute("method") || "post").toLowerCase();
            const data = formMethod !== "get" ? Tools.serializeFormDataToJSON(formData) : Tools.serializeFormData(formData);
            const buttons = formElement.querySelectorAll("button:not([disabled])");
            const toggleButtons = (disable) => {
              buttons.forEach((button) => button.disabled = disable);
            };
            toggleButtons(true);
            const action = formElement.getAttribute("action");
            if (!action) {
              toggleButtons(false);
              console.warn("Missing form action attribute. Skipping API call.");
              return;
            }
            API.makeRequest(
              formMethod,
              Tools.getBaseURL(action),
              data,
              (result) => {
                toggleButtons(false);
                API._afterComplete(formElement, result);
                return result;
              },
              (error) => {
                toggleButtons(false);
                FOSSBilling.ui.notify(`${error.message} (${error.code})`, "error");
              }
            );
          });
        });
      }
    },
    /**
     * Attach event listeners to links with data attribute 'data-fb-api'.
     **/
    _apiLink: function() {
      const linkElements = document.querySelectorAll("a[data-fb-api]");
      if (linkElements.length > 0) {
        linkElements.forEach((linkElement) => {
          if (linkElement.dataset.fbApiBound === "true") {
            return;
          }
          linkElement.dataset.fbApiBound = "true";
          const loadingState = createLinkLoadingState(linkElement);
          linkElement.addEventListener("click", function(event) {
            event.preventDefault();
            if (loadingState.isInProgress()) {
              return;
            }
            if (linkElement.classList.contains("disabled") || linkElement.getAttribute("aria-disabled") === "true") {
              return;
            }
            let apiData;
            try {
              apiData = Tools.parseDataAttr(linkElement.dataset.fbApi || "{}");
            } catch (error) {
              console.error("Failed to parse data-fb-api attribute:", error);
              FOSSBilling.ui.notify("Invalid API configuration", "error");
              return;
            }
            const rawHref = linkElement.getAttribute("href") || "";
            if (!apiData.href && (!rawHref || rawHref === "#")) {
              return;
            }
            const handleApiRequest = (method, href, params = {}) => {
              var _a, _b;
              if (apiData.loading || apiData.preventNavigation) {
                loadingState.set(apiData);
              }
              const url = apiData.href || href;
              const mergedParams = apiData.params && typeof apiData.params === "object" ? Object.assign({}, apiData.params, params) : params;
              API.makeRequest(
                method,
                Tools.getBaseURL(url),
                mergedParams,
                (result) => {
                  loadingState.reset();
                  API._afterComplete(linkElement, result);
                },
                (error) => {
                  loadingState.reset();
                  FOSSBilling.ui.notify(`${error.message} (${error.code})`, "error");
                },
                true,
                (_a = apiData.timeoutMs) != null ? _a : 3e4,
                (_b = apiData.timeoutMessage) != null ? _b : null
              );
            };
            const modalsLib = typeof Modals !== "undefined" ? Modals : null;
            dispatchLinkAction(apiData, rawHref, modalsLib, handleApiRequest);
          });
        });
      }
    }
  };
  window.FOSSBilling = window.FOSSBilling || {};
  window.FOSSBilling.tools = Tools;
  window.FOSSBilling.api = API;
  var bindApiInteractions = () => {
    if (document.querySelector("form[data-fb-api]")) {
      API._apiForm();
    }
    if (document.querySelector("a[data-fb-api]")) {
      API._apiLink();
    }
  };
  if (typeof FOSSBilling.ready === "function") {
    FOSSBilling.ready(bindApiInteractions);
  } else if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", bindApiInteractions);
  } else {
    bindApiInteractions();
  }
})();
//# sourceMappingURL=api.js.map
