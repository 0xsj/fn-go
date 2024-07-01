import viteSVG from "@/assets/vite-logo.svg";
export const ViteLogo: React.FC = () => {
  return (
    <img
      className='relative m-auto'
      src={viteSVG}
      width={301}
      height={60}
      alt='Vite Logo'
    />
  );
};
