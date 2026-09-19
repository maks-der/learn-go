(function () {
  var root = document.documentElement;
  var button = document.getElementById("theme-toggle");
  if (button) {
    var themes = ["light", "dark", "contrast-light", "contrast-dark"];
    var labels = {
      light: button.getAttribute("data-label-light") || "Light",
      dark: button.getAttribute("data-label-dark") || "Dark",
      "contrast-light": button.getAttribute("data-label-contrast-light") || "HC Light",
      "contrast-dark": button.getAttribute("data-label-contrast-dark") || "HC Dark",
    };

    function currentTheme() {
      var t = root.getAttribute("data-theme");
      return themes.indexOf(t) >= 0 ? t : "light";
    }

    function nextTheme(theme) {
      return themes[(themes.indexOf(theme) + 1) % themes.length];
    }

    function apply(theme) {
      root.setAttribute("data-theme", theme);
      try {
        localStorage.setItem("learngo-theme", theme);
      } catch (e) {}
      var next = nextTheme(theme);
      var label = button.querySelector("[data-theme-label]");
      if (label) {
        label.textContent = labels[next];
      }
      button.setAttribute("aria-label", labels[next]);
    }

    apply(currentTheme());
    button.addEventListener("click", function () {
      apply(nextTheme(currentTheme()));
    });
  }

  var header = document.querySelector(".top");
  if (!header) {
    return;
  }

  var reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
  var lastY = window.scrollY || 0;
  var ticking = false;

  function onScroll() {
    if (reduceMotion.matches) {
      header.classList.remove("is-hidden");
      return;
    }
    var y = window.scrollY || 0;
    if (y < 16) {
      header.classList.remove("is-hidden");
    } else if (y > lastY + 6) {
      header.classList.add("is-hidden");
    } else if (y < lastY - 6) {
      header.classList.remove("is-hidden");
    }
    lastY = y;
  }

  window.addEventListener(
    "scroll",
    function () {
      if (ticking) {
        return;
      }
      ticking = true;
      window.requestAnimationFrame(function () {
        onScroll();
        ticking = false;
      });
    },
    { passive: true }
  );
})();
