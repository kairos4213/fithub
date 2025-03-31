const token = sessionStorage.getItem("token");
const authSections = document.getElementsByClassName("auth-only");
const landing = document.getElementById("landing");

if (token) {
  for (let section of authSections) {
    section.style.display = "block";
  }
  landing.remove();
}
