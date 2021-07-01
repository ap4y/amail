module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      gridTemplateColumns: {
        layout: "14rem 1fr",
      },
      gridTemplateRows: {
        layout: "3.5rem 1fr",
      },
    },
  },
  variants: {
    extend: {
      textColor: ["visited", "active"],
      fontWeight: ["hover", "active"],
      backgroundColor: ["active", "disabled"],
      opacity: ["disabled"],
    },
  },
  plugins: [],
};
