(function () {
  function applyTheme() {
    var savedTheme = localStorage.getItem("theme") || "fithub-dark";
    document.documentElement.setAttribute("data-theme", savedTheme);

    var modeToggle = document.getElementById("theme-mode-toggle");
    if (modeToggle) {
      modeToggle.checked = savedTheme === "fithub-light";
      modeToggle.onchange = function (e) {
        var newTheme = e.target.checked ? "fithub-light" : "fithub-dark";
        localStorage.setItem("theme", newTheme);
        document.documentElement.setAttribute("data-theme", newTheme);
      };
    }

    var themeToggle = document.getElementById("theme-toggle");
    if (themeToggle) themeToggle.style.visibility = "visible";
  }

  // Apply theme immediately to prevent flash
  var savedTheme = localStorage.getItem("theme") || "fithub-dark";
  document.documentElement.setAttribute("data-theme", savedTheme);

  // Initialize on first load
  document.addEventListener("DOMContentLoaded", applyTheme);

  // Re-apply after htmx boosted navigation
  document.addEventListener("htmx:afterSwap", applyTheme);
})();
