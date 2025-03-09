document.addEventListener("DOMContentLoaded", async () => {
  const token = localStorage.getItem("token");

  if (token) {
    document.getElementById("auth-section").style.display = "none";
    document.getElementById("user-home").style.display = "block";
    await getWorkouts();
  } else {
    document.getElementById("auth-section").style.display = "block";
    document.getElementById("user-home").style.display = "none";
  }
});

document.getElementById("login-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  await login();
});

async function login() {
  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;

  try {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password}),
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`Failed to login: ${data.error}`);
    }
    
    console.log(data);
    if (data.access_token) {
      localStorage.setItem("token", data.access_token);
      document.getElementById("auth-section").style.display = "none";
      document.getElementById("user-home").style.display = "block";
      await getWorkouts();
    } else {
      alert("Login failed. Please check credentials.");
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
}
