document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('retrieveAuthorizedUsersButton').addEventListener('click', function() {
        const requestData = {
            endpointVar: 'getUsers',
        };

        fetch('https://admin.ryanschnabel.com/bff/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.text())
        .then(data => {
            const usersArray = data.split(',');

            // Get the div element
            const usersDiv = document.getElementById('whitelistedUsers');

            // Clear the current content
            usersDiv.innerHTML = '';

            // Create a span for each user, append it to the div, and add a line break
            usersArray.forEach(user => {
                const userSpan = document.createElement('span');
                userSpan.textContent = user.trim();
                usersDiv.appendChild(userSpan);
                usersDiv.appendChild(document.createElement('br'));
            });
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});
