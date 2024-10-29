// Add this to a new file: static/js/search.js
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const suggestionsContainer = document.getElementById('searchSuggestions');
    let debounceTimer;

    // Show/hide suggestions container
    searchInput.addEventListener('focus', () => {
        if (searchInput.value.length >= 2) {
            suggestionsContainer.style.display = 'block';
        }
    });

    document.addEventListener('click', (e) => {
        if (!searchInput.contains(e.target) && !suggestionsContainer.contains(e.target)) {
            suggestionsContainer.style.display = 'none';
        }
    });

    // Handle input changes
    searchInput.addEventListener('input', function() {
        clearTimeout(debounceTimer);
        const query = this.value.trim();
        
        if (query.length < 2) {
            suggestionsContainer.style.display = 'none';
            return;
        }

        debounceTimer = setTimeout(() => {
            fetchSuggestions(query);
        }, 300);
    });

    async function fetchSuggestions(query) {
        try {
            const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
            const data = await response.json();
            
            if (data.success) {
                displaySuggestions(data.results);
            }
        } catch (error) {
            console.error('Error fetching suggestions:', error);
        }
    }

    function displaySuggestions(results) {
        if (!results.length) {
            suggestionsContainer.style.display = 'none';
            return;
        }

        suggestionsContainer.innerHTML = '';
        
        results.forEach(result => {
            const div = document.createElement('div');
            div.className = 'suggestion-item';
            
            div.innerHTML = `
                <span class="suggestion-type">${result.type}</span>
                <span>${result.text}</span>
            `;
            
            div.addEventListener('click', () => {
                if (result.type === 'artist') {
                    window.location.href = `/artist?id=${result.id}`;
                } else if (result.type === 'location') {
                    window.location.href = `/locations?id=${result.id}`;
                } else if (result.type === 'date') {
                    window.location.href = `/dates?id=${result.id}`;
                }
            });

            suggestionsContainer.appendChild(div);
        });
        
        suggestionsContainer.style.display = 'block';
    }
});