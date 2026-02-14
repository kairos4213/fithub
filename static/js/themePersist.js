(function () {
  function applyTheme() {
    const savedTheme = localStorage.getItem("theme") || "fithub-dark";
    document.querySelector("html").setAttribute("data-theme", savedTheme);

    const themeToggle = document.getElementById("theme-toggle");
    const themeController = document.querySelector(".theme-controller");
    if (themeController) {
      themeController.checked = savedTheme === "fithub-light";
      themeController.addEventListener("change", (e) => {
        const newTheme = e.target.checked ? "fithub-light" : "fithub-dark";
        localStorage.setItem("theme", newTheme);
        document.querySelector("html").setAttribute("data-theme", newTheme);
      });
    }
    if (themeToggle) themeToggle.style.visibility = "visible";
  }

  // Apply theme immediately to prevent flash
  const savedTheme = localStorage.getItem("theme") || "fithub-dark";
  document.querySelector("html").setAttribute("data-theme", savedTheme);

  // Initialize on first load
  document.addEventListener("DOMContentLoaded", applyTheme);

  // Re-apply after htmx boosted navigation (afterSwap fires before browser paint)
  document.addEventListener("htmx:afterSwap", applyTheme);
})();
