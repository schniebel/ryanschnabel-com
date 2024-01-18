document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('retrieveAuthorizedUsersButton').addEventListener('click', function() {
        // ... existing fetch for getUsers ...

        .then(data => {
            const usersArray = data.split(',');

            const usersDiv = document.getElementById('whitelistedUsers');
            usersDiv.innerHTML = '';

            usersArray.forEach(user => {
                const userSpan = document.createElement('span');
                userSpan.textContent = user.trim();
                usersDiv.appendChild(userSpan);

                // Create Remove User Button
                const removeButton = document.createElement('button');
                removeButton.textContent = 'Remove User';
                removeButton.onclick = function() { removeUser(user.trim()); };
                usersDiv.appendChild(removeButton);

                usersDiv.appendChild(document.createElement('br'));
            });
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
    document.getElementById('addUserButton').addEventListener('click', function() {
        const inputText = document.getElementById('inputText').value;
        const requestData = {
            endpointVar: 'addUser',
            inputText: inputText
        };

        fetch('https://admin.ryanschnabel.com/bff/', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.text())
        .then(data => {
            // Process the response data
            console.log(data);
            // You can also update the UI here based on the response
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});

function removeUser(userEmail) {
    // Implement the logic to remove the user
    // For example, you might send a DELETE request to your backend
    console.log('Remove user:', userEmail);
    // Add here the code to call your backend to remove the user
}