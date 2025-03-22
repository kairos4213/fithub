const Router = {
  init: () => {
    document.querySelectorAll("a.navlink").forEach((a) => {
      a.addEventListener("click", (event) => {
        event.preventDefault();
        const href = event.target.getAttribute("href");
        Router.go(href);
      });
    });

    window.addEventListener("popstate", (event) => {
      Router.go(event.state.route, false);
    });

    // Process initial URL
    Router.go(location.pathname);
  },
  go: (route, addToHistory = true) => {
    if (addToHistory) {
      history.pushState({ route }, "", route);
    }
    let pageElement = null;
    switch (route) {
      case "/app/":
        pageElement = document.createElement("h2");
        pageElement.textContent = "Home";
        break;
      case "/app/login":
        pageElement = document.createElement("h2");
        pageElement.textContent = "Login";
        break;
      case "/app/goals/":
        pageElement = document.createElement("h2");
        pageElement.textContent = "Goals";
        pageElement.dataset.goalId = route.substring(
          route.lastIndexOf("/") + 1,
        );
        break;
      default:
        pageElement = document.createElement("h2");
        pageElement.textContent = "Where are you headed?";
    }
    if (pageElement) {
      document.querySelector("main").innerHTML = "";
      document.querySelector("main").appendChild(pageElement);
    }

    window.scrollX = 0;
  },
};

export default Router;
