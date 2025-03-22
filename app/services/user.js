import API from "./api.js";

export async function loginUser() {
  const data = await API.login();
  app.store.user = data;
}
