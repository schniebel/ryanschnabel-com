document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('callApiButton').addEventListener('click', function() {
        var inputText = document.getElementById('inputText').value;
        
        fetch('https://api.ryanschnabel.com/helloworld')
            .then(response => response.text())
            .then(data => {
                var combinedResponse = data + ' ' + inputText;
                document.getElementById('apiResponse').textContent = combinedResponse;
            })
            .catch(error => {
                console.error('Error:', error);
                document.getElementById('apiResponse').textContent = 'Error: ' + error.message;
            });
    });
});