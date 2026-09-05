(() => {
  // frontend/core/currency-format.mts
  function getFractionDigits(options) {
    const value = options.fraction_digits;
    return Number.isInteger(value) && value >= 0 && value <= 6 ? value : void 0;
  }
  function formatCurrencyAmount(amount, options, locale) {
    var _a;
    const fractionDigits = getFractionDigits(options);
    const fractionOptions = fractionDigits === void 0 ? {} : {
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits
    };
    const pattern = (_a = options.format_pattern) == null ? void 0 : _a.trim();
    if (!pattern || pattern.split("{amount}").length !== 2) {
      return new Intl.NumberFormat(locale, {
        style: "currency",
        currency: options.code,
        currencyDisplay: "narrowSymbol",
        ...fractionOptions
      }).format(amount);
    }
    let decimalFractionOptions = fractionOptions;
    if (fractionDigits === void 0) {
      const currencyOptions = new Intl.NumberFormat(locale, {
        style: "currency",
        currency: options.code
      }).resolvedOptions();
      decimalFractionOptions = {
        minimumFractionDigits: currencyOptions.minimumFractionDigits,
        maximumFractionDigits: currencyOptions.maximumFractionDigits
      };
    }
    const formattedAmount = new Intl.NumberFormat(locale, {
      style: "decimal",
      ...decimalFractionOptions
    }).format(amount);
    return pattern.replace("{amount}", formattedAmount);
  }

  // frontend/core/fossbilling.ts
  (function(window2, document2) {
    "use strict";
    const FOSSBilling = window2.FOSSBilling || {};
    const readyCallbacks = [];
    const editorsByElement = /* @__PURE__ */ new WeakMap();
    const editorsByName = /* @__PURE__ */ new Map();
    const adapters = /* @__PURE__ */ new Map();
    const cookieNames = Object.freeze({
      locale: "fossbilling_locale",
      timezone: "fossbilling_timezone"
    });
    function runReadyCallback(callback) {
      try {
        callback(FOSSBilling);
      } catch (error) {
        console.error("FOSSBilling ready callback failed:", error);
      }
    }
    FOSSBilling.ready = function(callback) {
      if (typeof callback !== "function") {
        return;
      }
      if (document2.readyState === "loading") {
        readyCallbacks.push(callback);
        return;
      }
      runReadyCallback(callback);
    };
    document2.addEventListener("DOMContentLoaded", function() {
      while (readyCallbacks.length > 0) {
        runReadyCallback(readyCallbacks.shift());
      }
      migrateCookie(cookieNames.locale, "fb_locale", 365);
      initTimezone();
    });
    FOSSBilling.ui = Object.assign({
      notify(message, type = "info") {
        if (typeof FOSSBilling.message === "function") {
          FOSSBilling.message(message, type);
          return;
        }
        if (type === "error") {
          console.error(message);
          return;
        }
        console.info(message);
      }
    }, FOSSBilling.ui || {});
    FOSSBilling.cookieCreate = FOSSBilling.cookieCreate || function(name, value, days) {
      let expires = "";
      if (days) {
        const date = /* @__PURE__ */ new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1e3);
        expires = `; expires=${date.toUTCString()}`;
      }
      document2.cookie = `${name}=${value}${expires}; path=/`;
    };
    FOSSBilling.cookieRead = FOSSBilling.cookieRead || function(name) {
      const nameEQ = `${name}=`;
      const cookies = document2.cookie.split(";");
      for (let i = 0; i < cookies.length; i++) {
        let cookie = cookies[i];
        while (cookie.charAt(0) === " ") {
          cookie = cookie.substring(1);
        }
        if (cookie.indexOf(nameEQ) === 0) {
          return cookie.substring(nameEQ.length);
        }
      }
      return null;
    };
    function migrateCookie(name, legacyName, days) {
      const currentValue = FOSSBilling.cookieRead(name);
      const legacyValue = FOSSBilling.cookieRead(legacyName);
      if (!currentValue && legacyValue) {
        FOSSBilling.cookieCreate(name, legacyValue, days);
      }
      if (legacyValue !== null) {
        FOSSBilling.cookieCreate(legacyName, "", -1);
      }
      return currentValue || legacyValue;
    }
    FOSSBilling.cookieNames = cookieNames;
    FOSSBilling.currency = Object.assign({
      format: formatCurrencyAmount
    }, FOSSBilling.currency || {});
    FOSSBilling.detectTimezone = FOSSBilling.detectTimezone || function() {
      try {
        const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
        if (typeof tz === "string" && tz.length > 0) {
          return tz;
        }
      } catch (error) {
      }
      return null;
    };
    function initTimezone() {
      const detected = FOSSBilling.detectTimezone();
      if (!detected) {
        return;
      }
      if (!migrateCookie(cookieNames.timezone, "fb_timezone", 365)) {
        FOSSBilling.cookieCreate(cookieNames.timezone, detected, 365);
      }
      const selects = document2.querySelectorAll("select[data-timezone-select]");
      selects.forEach(function(select) {
        if (select.value) {
          return;
        }
        const option = Array.from(select.options).find(function(opt) {
          return opt.value === detected;
        });
        if (option) {
          select.value = detected;
        }
      });
    }
    FOSSBilling.initTimezone = initTimezone;
    function getEditorName(element) {
      return element.getAttribute("name") || element.id || null;
    }
    function syncElementData(element, editor) {
      const data = editor.getData();
      if ("value" in element) {
        element.value = data;
      } else {
        element.textContent = data;
      }
    }
    function normalizeEditor(element, editor) {
      if (!editor || typeof editor.getData !== "function" || typeof editor.setData !== "function") {
        throw new Error("Editor adapters must return getData() and setData() methods.");
      }
      return {
        raw: editor.raw || editor,
        getData: () => editor.getData(),
        setData: (value) => editor.setData(value),
        focus: () => {
          var _a, _b;
          if (typeof editor.focus === "function") {
            editor.focus();
            return;
          }
          if ((_b = (_a = editor.editing) == null ? void 0 : _a.view) == null ? void 0 : _b.focus) {
            editor.editing.view.focus();
          }
        },
        destroy: () => {
          if (typeof editor.destroy === "function") {
            return editor.destroy();
          }
          return Promise.resolve();
        },
        onChange: (callback) => {
          var _a, _b;
          if (typeof editor.onChange === "function") {
            return editor.onChange(callback);
          }
          if ((_b = (_a = editor.model) == null ? void 0 : _a.document) == null ? void 0 : _b.on) {
            editor.model.document.on("change:data", callback);
          }
        },
        element,
        name: getEditorName(element),
        required: element.dataset.editorRequired === "true"
      };
    }
    FOSSBilling.editor = Object.assign({
      registerAdapter(name, adapter) {
        if (!name || !adapter || typeof adapter.create !== "function") {
          throw new Error("Editor adapter registration requires a name and create() method.");
        }
        adapters.set(name, adapter);
      },
      async create(element, options = {}) {
        if (!element) {
          throw new Error("Cannot initialize an editor without an element.");
        }
        const adapterName = options.adapter || "ckeditor";
        const adapter = adapters.get(adapterName);
        if (!adapter) {
          throw new Error(`Editor adapter "${adapterName}" is not registered.`);
        }
        if (element.hasAttribute("required")) {
          element.dataset.editorRequired = "true";
          element.removeAttribute("required");
        }
        const rawEditor = await adapter.create(element, options);
        const editor = normalizeEditor(element, rawEditor);
        syncElementData(element, editor);
        editor.onChange(() => syncElementData(element, editor));
        editorsByElement.set(element, editor);
        if (editor.name) {
          editorsByName.set(editor.name, editor);
        }
        element.editor = editor;
        element.setAttribute("data-editor", "true");
        return editor;
      },
      init(selector, options = {}) {
        document2.querySelectorAll(selector).forEach((element) => {
          if (editorsByElement.has(element)) {
            return;
          }
          FOSSBilling.editor.create(element, options).catch((error) => {
            console.error("Editor initialization error:", error);
          });
        });
      },
      get(name) {
        return editorsByName.get(name) || null;
      },
      getForElement(element) {
        return editorsByElement.get(element) || null;
      },
      all() {
        return Array.from(editorsByName.values());
      },
      syncForm(form, formData) {
        const formElements = new Set(Array.from(form.elements || []));
        editorsByName.forEach((editor, name) => {
          if (!formElements.has(editor.element)) {
            return;
          }
          syncElementData(editor.element, editor);
          formData.set(name, editor.getData());
        });
      },
      validateForm(form) {
        const formElements = new Set(Array.from(form.elements || []));
        for (const editor of editorsByName.values()) {
          if (!formElements.has(editor.element)) {
            continue;
          }
          if (editor.required && editor.getData().trim() === "") {
            return false;
          }
        }
        return true;
      }
    }, FOSSBilling.editor || {});
    window2.FOSSBilling = FOSSBilling;
  })(window, document);
})();
//# sourceMappingURL=fossbilling.js.map
