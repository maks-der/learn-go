(function () {
  var root = document.documentElement;
  var button = document.getElementById("theme-toggle");
  if (!button) {
    return;
  }

  var labels = {
    light: button.getAttribute("data-label-light") || "Light",
    dark: button.getAttribute("data-label-dark") || "Dark",
  };

  function currentTheme() {
    return root.getAttribute("data-theme") === "dark" ? "dark" : "light";
  }

  function apply(theme) {
    root.setAttribute("data-theme", theme);
    try {
      localStorage.setItem("learngo-theme", theme);
    } catch (e) {}
    var next = theme === "dark" ? "light" : "dark";
    var label = button.querySelector("[data-theme-label]");
    if (label) {
      label.textContent = next === "dark" ? labels.dark : labels.light;
    }
    button.setAttribute("aria-pressed", theme === "dark" ? "true" : "false");
  }

  apply(currentTheme());
  button.addEventListener("click", function () {
    apply(currentTheme() === "dark" ? "light" : "dark");
  });
})();
