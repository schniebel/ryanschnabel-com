document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('callApiButton').addEventListener('click', function() {
        // The data to send in the request body
        const requestData = {
            endpointVar: 'helloworld'
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
            // Handle the response data
            console.log(data);
        })
        .catch(error => {
            // Handle any errors
            console.error('Error:', error);
        });
    });
});
