// src/hooks/usei18n.ts
import { useTranslation } from "react-i18next";

type UseI18nResult = {
  changeLanguage: (lng: string) => void;
};

const useI18n = (): UseI18nResult => {
  const { i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return { changeLanguage };
};

export default useI18n;
