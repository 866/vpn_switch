async function sendRequest() {
    const response = await fetch("https://jsonplaceholder.typicode.com/posts/1");
    const post = await response.json();
    
    // Для наглядності, що запит відправляється
    document.getElementById('res').innerText = JSON.stringify(post);
  }
  
  
  const switchEl = document.getElementById('switch');
  
  switchEl.addEventListener('change', (event) => {
    if (event.currentTarget.checked) {
      console.log('checked');
      sendRequest();
    } else {
      console.log('not checked');
    }
  });
  