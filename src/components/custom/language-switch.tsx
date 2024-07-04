import React, { useCallback } from "react";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import useI18n from "@/hooks/use-i18n";

export const LanguageSwitch: React.FC = () => {
  const { changeLanguage } = useI18n();

  const handleLanguage = (value: string) => {
    changeLanguage(value.toLowerCase());
  };
  return (
    <Select onValueChange={handleLanguage}>
      <SelectTrigger className='w-[80px]'>
        <SelectValue placeholder='EN' />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectItem value='EN'>EN</SelectItem>
          <SelectItem value='KR'>KR</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  );
};
