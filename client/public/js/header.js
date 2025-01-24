document.addEventListener("DOMContentLoaded", () => {
    const menu = document.querySelector(".primary-menu");
    
    // Tema kontrolü için bir fonksiyon
    function updateTheme() {
      const bgColor = getComputedStyle(menu).backgroundColor;
      if (bgColor === "rgb(255, 255, 255)") {
        menu.setAttribute("data-bs-theme", "light");
      } else {
        menu.setAttribute("data-bs-theme", "dark");
      }
    }
  
    // Tema kontrolü çağrısı
    updateTheme();
  
    // Pencere boyutu değiştiğinde kontrol et
    window.addEventListener("resize", updateTheme);
  });