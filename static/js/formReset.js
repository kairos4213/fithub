/**
 * Resets a form and clears any associated field error elements.
 * @param {string} formID - The ID of the form element to reset.
 * @param {string[]} errorIDs - Array of error element IDs to hide and clear.
 */
function resetForm(formID, errorIDs) {
  const form = document.getElementById(formID);
  if (form && typeof form.reset === "function") {
    form.reset();
  }
  for (const id of errorIDs) {
    const el = document.getElementById(id);
    if (el) {
      el.className = "hidden";
      el.innerHTML = "";
    }
  }
}

/**
 * After HTMX settles OOB swaps, find the first visible field error
 * and scroll/focus its associated input.
 */
function scrollToFirstError() {
  const errEls = document.querySelectorAll("div[id^='err-']:not(.hidden)");
  if (!errEls.length) return;

  const errEl = errEls[0];
  let input = errEl.previousElementSibling;
  if (input && !input.matches("input, textarea, select")) {
    input = input.querySelector("input, textarea, select");
  }

  errEl.scrollIntoView({ behavior: "smooth", block: "center" });
  if (input) {
    setTimeout(() => input.focus(), 300);
  }
}

document.addEventListener("htmx:afterSettle", function (evt) {
  // Only run on 400 responses (validation errors)
  if (evt.detail.xhr && evt.detail.xhr.status === 400) {
    scrollToFirstError();
  }
});
