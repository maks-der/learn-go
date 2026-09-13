(function () {
  var root = document.documentElement;
  var button = document.getElementById("theme-toggle");
  if (button) {
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
