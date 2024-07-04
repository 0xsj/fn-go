import i18next from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locale/en/translations.json";
import es from "@/locale/es/translations.json";
import kr from "@/locale/kr/translations.json";

const resources = {
  en: {
    translation: en,
  },
  es: {
    translation: es,
  },
  kr: {
    translation: kr,
  },
};

i18next.use(initReactI18next).init({
  resources,
  lng: "en",
});

export default i18next;
