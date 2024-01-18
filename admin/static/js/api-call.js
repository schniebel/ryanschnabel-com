document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('retrieveAuthorizedUsersButton').addEventListener('click', function() {

        const requestData = {
            endpointVar: 'getUsers',
        };

        // Making a POST request to the BFF
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

            // Get the list element
            const userList = document.getElementById('whitelistedUsers');

            // Clear the current list
            userList.innerHTML = '';

            // Create a list item for each user and append it to the list
            usersArray.forEach(user => {
                const listItem = document.createElement('li');
                listItem.textContent = user.trim();
                userList.appendChild(listItem);
            });
        })
        .catch(error => {
            // Handle any errors
            console.error('Error:', error);
        });
    });
});
