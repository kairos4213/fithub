/**
 * Resets a form and clears any associated field error elements.
 * @param {string} formID - The ID of the form element to reset.
 * @param {string[]} errorIDs - Array of error element IDs to hide and clear.
 */
function resetForm(formID, errorIDs) {
  const form = document.getElementById(formID);
  if (form) {
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
