import { LanguageSwitch } from "@/components/custom/language-switch";
import { Layout } from "@/components/custom/layout";
import { SearchBox } from "@/components/custom/search-box";
import { ThemeSwitch } from "@/components/custom/theme-switch";
import { UserNav } from "@/components/user-nav";
import { useState } from "react";

interface Props {}

export const Documents: React.FC<Props> = () => {
  return (
    <Layout>
      <Layout.Header>
        <div className='ml-auto flex items-center space-x-4'>
          <SearchBox />
          <ThemeSwitch />
          <UserNav />
          <LanguageSwitch />
        </div>
      </Layout.Header>
      <Layout.Body>
        <div></div>
      </Layout.Body>
    </Layout>
  );
};

export default Documents;
