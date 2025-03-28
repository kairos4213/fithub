const token = sessionStorage.getItem("token");
const mainSection = document.getElementsByTagName("main")[0];

if (!token) {
  window.location.href = "../index.html";
  alert(
    "You're not authorized to view this, please login or register to continue",
  );
} else {
  mainSection.style.display = "block";
}
