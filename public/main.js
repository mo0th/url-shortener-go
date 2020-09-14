const form = document.getElementById("form");
const shortInput = document.getElementById("short-input");
const urlInput = document.getElementById("url-input");
const errDiv = document.getElementById("error");
const sucDiv = document.getElementById("success");
const resInput = document.getElementById("result");
const tryLink = document.getElementById("try");

const API_URL = "https://mighty-gorge-93392.herokuapp.com/";

urlInput.focus();

function showSuccessStuff(isSuccess) {
  sucDiv.classList[isSuccess ? "remove" : "add"]("hidden");
  resInput.classList[isSuccess ? "remove" : "add"]("hidden");
  tryLink.classList[isSuccess ? "remove" : "add"]("hidden");
  errDiv.classList[!isSuccess ? "remove" : "add"]("hidden");
}

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const res = await fetch(API_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      short: shortInput.value,
      url: urlInput.value,
    }),
  });
  const data = await res.json();

  if (data.success) {
    showSuccessStuff(true);
    sucDiv.textContent = "shortened";
    const link = `${API_URL}${data.short}`;
    resInput.value = link;
    tryLink.href = link;
  } else {
    showSuccessStuff(false);
    errDiv.textContent = data.message;
  }
});
