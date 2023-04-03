const generateRandomNumber = () => {
  const min = 100000000; // số nhỏ nhất có 9 chữ số
  const max = 999999999; // số lớn nhất có 9 chữ số
  const randomNumber = Math.floor(Math.random() * (max - min + 1) + min);
  return randomNumber.toString();
};

export default generateRandomNumber;
