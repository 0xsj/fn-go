import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

interface Props {}
export const SearchBox: React.FC<Props> = () => {
  const { t, i18n } = useTranslation();
  return (
    <div>
      <Input
        type='search'
        placeholder={`${t("search")}`}
        className={cn(`md:w-[100px] lg:w-[300px]`)}
      />
    </div>
  );
};
