document.addEventListener('DOMContentLoaded', function() {
    const readMoreButton = document.getElementById('read-more-button');
    let offset = 24;

    if (readMoreButton) {
        readMoreButton.addEventListener('click', function() {
            const urlParams = new URLSearchParams(window.location.search);
            const category = urlParams.get('category') || '';

            fetch(`/products?offset=${offset}&limit=24&category=${category}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                }
            })
            .then(response => response.json())
            .then(data => {
                const shopDiv = document.getElementById("shop");
                if (data.items && data.items.length > 0) {
                    const fragment = document.createDocumentFragment();
                    data.items.forEach(item => {
                        // ... mevcut ürün oluşturma kodu ...
                    });
                    shopDiv.appendChild(fragment);
                    offset += 24;
                    
                    // Eğer daha fazla ürün yoksa butonu gizle
                    if (!data.hasMore) {
                        readMoreButton.style.display = 'none';
                    }
                } else {
                    readMoreButton.style.display = 'none';
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert("Ürünler yüklenirken bir hata oluştu.");
            });
        });
    }
});