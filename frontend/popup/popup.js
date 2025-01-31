// Пример данных для списка криптовалют
const cryptoData = [
  {
    name: "Bitcoin",
    symbol: "BTC",
    logo: "https://cryptologos.cc/logos/bitcoin-btc-logo.svg?v=023",
    amount: 0.01912343,
    priceUsd: 0,
    priceChange: ""
  },
  {
    name: "Litecoin",
    symbol: "LTC",
    logo: "https://cryptologos.cc/logos/litecoin-ltc-logo.svg?v=023",
    amount: 20.00085,
    priceUsd: 0,
    priceChange: ""
  },
  {
    name: "TRON",
    symbol: "TRX",
    logo: "https://cryptologos.cc/logos/tron-trx-logo.svg?v=023",
    amount: 1580.8565,
    priceUsd: 0,
    priceChange: ""
  },
  {
    name: "Ethereum",
    symbol: "ETH",
    logo: "https://cryptologos.cc/logos/ethereum-eth-logo.svg?v=023",
    amount: 0.15628,
    priceUsd: 0,
    priceChange: ""
  }
];

// Функция для форматирования чисел
function formatNumber(num, decimals = 2) {
  return Number(num).toFixed(decimals);
}

// Суммируем все монеты * их цену
function calculateTotalBalance(data) {
  let total = 0;
  data.forEach((coin) => {
    total += coin.amount * coin.priceUsd;
  });
  return total; // сумма в долларах
}

// Отрисовка списка криптоактивов, как раньше
function renderCryptoList() {
  const listContainer = document.getElementById("cryptoList");
  listContainer.innerHTML = "";

  cryptoData.forEach((coin) => {
    const card = document.createElement("div");
    card.classList.add("crypto-card");

    const infoDiv = document.createElement("div");
    infoDiv.classList.add("crypto-info");

    const logoDiv = document.createElement("div");
    logoDiv.classList.add("crypto-logo");
    const logoImg = document.createElement("img");
    logoImg.src = coin.logo;
    logoImg.alt = coin.symbol;
    logoDiv.appendChild(logoImg);

    const nameDiv = document.createElement("div");
    const nameEl = document.createElement("div");
    nameEl.classList.add("crypto-name");
    nameEl.textContent = coin.name + " (" + coin.symbol + ")";
    nameDiv.appendChild(nameEl);

    infoDiv.appendChild(logoDiv);
    infoDiv.appendChild(nameDiv);

    const statsDiv = document.createElement("div");
    statsDiv.classList.add("crypto-stats");

    const amountEl = document.createElement("div");
    amountEl.classList.add("crypto-amount");
    amountEl.textContent = formatNumber(coin.amount, 8);

    const priceEl = document.createElement("div");
    priceEl.classList.add("crypto-price");
    priceEl.textContent = `$${formatNumber(coin.priceUsd)} (${coin.priceChange})`;

    statsDiv.appendChild(amountEl);
    statsDiv.appendChild(priceEl);

    card.appendChild(infoDiv);
    card.appendChild(statsDiv);
    listContainer.appendChild(card);
  });
}

document.addEventListener("DOMContentLoaded", () => {
  // Сначала рендерим список
  renderCryptoList();

  // Считаем общий баланс
  const total = calculateTotalBalance(cryptoData);

  // Получаем элементы с балансом
  const balanceSection = document.querySelector(".balance-section");
  const balanceElem = document.querySelector(".balance");
  const balanceDiffElem = document.querySelector(".balance-diff");

  // Если общий баланс > 0, показываем и записываем значение
  if (total > 0) {
    balanceElem.textContent = `$${formatNumber(total)}`;
    balanceDiffElem.textContent = "+ $242.54 (2.93%)"; 
    // Здесь можно вывести реальный процент изменения или другую логику 
    // (это зависит от того, как вы рассчитываете рост/падение)
  } else {
    // Иначе скрываем секцию целиком
    balanceSection.style.display = "none";
  }

  // Привязываем события к кнопкам
  document.getElementById("receiveBtn").addEventListener("click", () => {
    alert("Функция для получения монет");
  });

  document.getElementById("scanBtn").addEventListener("click", () => {
    alert("Функция для сканирования QR-кода");
  });

  document.getElementById("sendBtn").addEventListener("click", () => {
    alert("Функция для отправки монет");
  });
});
