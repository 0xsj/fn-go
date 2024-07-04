import React, { HTMLAttributes } from "react";
import { DropdownMenuCheckboxItemProps } from "@radix-ui/react-dropdown-menu";

import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";

// Create the context and its type
interface SwitchContextType {
  date: boolean;
  toggleDate: () => void;
}

const SwitchContext = React.createContext<SwitchContextType>({
  date: false,
  toggleDate: () => {},
});

// eslint-disable-next-line react-refresh/only-export-components
export const useSwitch = () => React.useContext(SwitchContext);

interface Props extends HTMLAttributes<HTMLInputElement> {}

export const SwitchProvider: React.FC<Props> = ({ children }) => {
  const [date, setDate] = React.useState<boolean>(false);

  const toggleDate = () => {
    setDate((prevDate) => !prevDate);
  };

  return (
    <SwitchContext.Provider value={{ date, toggleDate }}>
      {children}
    </SwitchContext.Provider>
  );
};

// Modify the GraphDataSwitch component to use the context
export const GraphDataSwitch: React.FC = () => {
  const { date, toggleDate } = useSwitch();

  return (
    <div className='flex items-center space-x-2'>
      <Switch id='date' checked={date} onCheckedChange={toggleDate} />
      <Label htmlFor='date'>View</Label>
    </div>
  );
};
