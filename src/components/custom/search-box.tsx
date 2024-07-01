import { Input } from "@/components/ui/Input"
import { cn } from "@/lib/utils"

interface Props {}
export const SearchBox: React.FC<Props> = () => {
  return (
    <div>
      <Input
        type='search'
        placeholder='Search...'
        className={cn(`md:w-[100px] lg:w-[300px]`)}
      />
    </div>
  )
}
