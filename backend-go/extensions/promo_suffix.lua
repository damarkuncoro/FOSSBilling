-- FOSSBilling Dynamic Extension Example
-- This script runs automatically via the Lua Engine hooks

fb.log("Promo Suffix extension initialized!")

-- Hook: filter_invoice_item_title
-- Appends a [LUA] tag to all invoice items
function filter_invoice_item_title(title)
    fb.log("Intercepting title: " .. title)
    return title .. " [Dinamis via Lua]"
end
