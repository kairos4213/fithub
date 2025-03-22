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
        if (sessionStorage.getItem("token") != null) {
          pageElement = document.createElement("home-page");
        } else {
          pageElement = document.createElement("landing-page");
        }
        break;
      case "/app/login":
        pageElement = document.createElement("login-page");
        break;
      case "/app/register":
        pageElement = document.createElement("register-page");
        break;
      case "/app/goals":
        pageElement = document.createElement("goals-page");
        break;
      case "/app/goals/":
        pageElement = document.createElement("goal-page");
        pageElement.dataset.goalId = route.substring(
          route.lastIndexOf("/") + 1,
        );
        break;
      default:
        pageElement = document.createElement("h2");
        pageElement.textContent = "This is not the page you were looking for";
    }
    if (pageElement) {
      document.querySelector("main").innerHTML = "";
      document.querySelector("main").appendChild(pageElement);
    }

    window.scrollX = 0;
  },
};

export default Router;
