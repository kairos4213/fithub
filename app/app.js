import Store from "./services/store.js";
import API from "./services/api.js";
import { loginUser } from "./services/user.js";
import Router from "./services/router.js";

window.app = {};
app.store = Store;
app.router = Router;

const $ = () => document.querySelector.call(this, arguments);
const $$ = () => document.querySelectorAll.call(this, arguments);
HTMLElement.prototype.on = (a, b, c) => this.addEventListener(a, b, c);
HTMLElement.prototype.off = (a, b) => this.removeEventListener(a, b);
HTMLElement.prototype.$ = (s) => this.querySelector(s);
HTMLElement.prototype.$ = (s) => this.querySelectorAll(s);

window.addEventListener("DOMContentLoaded", () => {
  app.router.init();
  loginUser();
});
