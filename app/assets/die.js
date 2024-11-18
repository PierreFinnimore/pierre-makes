const DICE_SEPARATOR = "-";

function getRolledDice() {
  const params = new URLSearchParams(window.location.search);
  const dieParam = params.get("dice");
  const maxVals = dieParam ? dieParam.split(DICE_SEPARATOR) : [];
  const isRolling = false;
  return maxVals.map((maxStr) => {
    const max = Number(maxStr);
    const value = getDieRoll(max);
    return { max, value, isRolling };
  });
}

function getDieRoll(max) {
  return Math.floor(Math.random() * max + 1);
}

document.addEventListener("alpine:init", () => {
  Alpine.data("dice", () => ({
    rolledDice: getRolledDice(),
    rollDice: function () {
      this.rolledDice = getRolledDice();
      this.addRollingAnimation(0);
    },
    rollDie: function (index) {
      this.rolledDice[index].value = getDieRoll(this.rolledDice[index].max);
      this.rolledDice[index].isRolling = true;
      setTimeout(() => {
        this.rolledDice[index].isRolling = false;
      }, 150);
    },
    sumDice: function () {
      return this.rolledDice.reduce((acc, curr) => {
        return acc + curr.value;
      }, 0);
    },
    meanDice: function () {
      return (
        this.rolledDice.reduce((acc, curr) => {
          return acc + curr.value;
        }, 0) / Math.max(this.rolledDice.length, 1)
      );
    },
    modeDice: function () {
      const valueCounts = {};
      this.rolledDice.forEach((die) => {
        if (die.value in valueCounts) {
          valueCounts[die.value] += 1;
        } else {
          valueCounts[die.value] = 1;
        }
      });
      let highestVal = 0;
      let highestKeys = [];

      Object.entries(valueCounts).forEach(([key, val]) => {
        if (val >= highestVal) {
          if (val === highestVal) {
            highestKeys.push(key);
          } else {
            highestKeys = [key];
          }
          highestVal = val;
        }
      });
      return highestKeys.join(", ");
    },
    addDice: function (max, count) {
      const params = new URLSearchParams(window.location.search);
      const dieParam = params.get("dice");
      const maxVals = dieParam ? dieParam.split(DICE_SEPARATOR) : [];

      const startIndex = maxVals.length;
      for (let i = 1; i <= count; i++) {
        maxVals.push(max);
        this.rolledDice.push({ max, value: getDieRoll(max), isRolling: true });
      }
      this.addRollingAnimation(startIndex);

      params.set("dice", maxVals.join(DICE_SEPARATOR));
      const currentUrl = new URL(window.location.href);
      currentUrl.search = params.toString();
      history.replaceState(null, "", currentUrl);
    },
    removeDie: function (index) {
      const params = new URLSearchParams(window.location.search);
      const dieParam = params.get("dice");
      const maxVals = dieParam ? dieParam.split(DICE_SEPARATOR) : [];
      maxVals.splice(index, 1);
      params.set("dice", maxVals.join(DICE_SEPARATOR));
      const currentUrl = new URL(window.location.href);
      currentUrl.search = params.toString();
      history.replaceState(null, "", currentUrl);
      this.rolledDice.splice(index, 1);
    },
    deleteAllDice: function () {
      this.rolledDice = [];
      const url = new URL(window.location);
      url.searchParams.delete("dice");
      history.replaceState({}, "", url);
    },
    addRollingAnimation(startIndex) {
      this.rolledDice.forEach((die, index) => {
        if (index >= startIndex) {
          die.isRolling = true;
          setTimeout(() => {
            die.isRolling = false;
          }, 150);
        }
      });
    },
    copyText: "Share",
    getWindowUrl: function () {
      navigator.clipboard.writeText(window.location.href);
      this.copyText = "Copied!";
      setTimeout(() => {
        this.copyText = "Share";
      }, 500);
    },
  }));
});
