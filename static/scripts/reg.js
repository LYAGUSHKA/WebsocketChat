document.querySelector('.reg__register').addEventListener('click', () => {
  document.querySelector('.reg__login').classList.remove('active');  
  document.querySelector('.reg__register').classList.add('active');
  document.querySelector('.login').classList.add('disabled')
  document.querySelector('.register').classList.remove('disabled') 
});


document.querySelector('.reg__login').addEventListener('click', () => {
  document.querySelector('.reg__register').classList.remove('active');  
  document.querySelector('.reg__login').classList.add('active');
  document.querySelector('.register').classList.add('disabled')
  document.querySelector('.login').classList.remove('disabled') 
});