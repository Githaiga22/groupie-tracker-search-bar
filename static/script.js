document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const suggestionsContainer = document.getElementById('searchSuggestions');
    let debounceTimer;

    // Show/hide suggestions container
    searchInput.addEventListener('focus', () => {
        if (searchInput.value.length > 0) {
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
        
        if (query.length === 0) {
            suggestionsContainer.style.display = 'none';
            return;
        }

        debounceTimer = setTimeout(() => {
            fetchSuggestions(query);
        }, 150); // Reduced debounce time
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
            // Format location before displaying
            function formatLocation(location) {
                return location
                    .replace(/[_]/g, ' ')  
                    .replace(/[-]/g, ' , ')              // Replace underscores and hyphens with space
                    .split(' ')                           // Split into words
                    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()) // Capitalize first letter of each word
                    .join(' ');                          // Join words back into a single string
            }
            
            // Display both the main text and context if available
            div.innerHTML = `
                <span class="suggestion-type">${result.type}</span>
                ${result.context ? `<span class="suggestion-type">${result.context}</span>` : ''}
                <span>${formatLocation(result.text)}</span>
                
            `;
            
            div.addEventListener('click', () => {
                    window.location.href = `/artist?id=${result.id}`
            });
    
            suggestionsContainer.appendChild(div);
        });
        
        suggestionsContainer.style.display = 'block';
    }
});
