const API = {
  url: "/api/v1/",
  login: async () => {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    try {
      const res = await fetch(API.url + "login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email,
          password,
        }),
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(`Failed to login: ${data.Error}`);
      }

      if (data.access_token) {
        sessionStorage.setItem("token", data.access_token);
      } else {
        alert("Login failed, please check credentials");
      }
      return data;
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  },
};

export default API;
