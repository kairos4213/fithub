document
  .getElementById("register-form")
  .addEventListener("submit", async (event) => {
    event.preventDefault();
    login();
  });

async function login() {
  const first_name = document.getElementById("first_name").value;
  const last_name = document.getElementById("last_name").value;
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;

  try {
    const res = await fetch("/api/v1/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        first_name,
        last_name,
        email,
        password,
      }),
    });

    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to register: ${data.error}`);
    }

    if (data.access_token) {
      sessionStorage.setItem("token", data.access_token);
      window.location.href = "../pages/home.html";
    } else {
      alert("Failed to register new user. Please try again.");
    }
    return data;
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}
