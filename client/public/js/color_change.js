document.addEventListener("DOMContentLoaded", function () {
    const header = document.getElementById("header");

    window.addEventListener("scroll", function () {
        if (window.scrollY > 50) {
            header.classList.add("sticky");
        } else {
            header.classList.remove("sticky");
        }
    });
});
