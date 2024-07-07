import { Button } from "@/components/custom/button";
import { LanguageSwitch } from "@/components/custom/language-switch";
import { Layout } from "@/components/custom/layout";
import { SearchBox } from "@/components/custom/search-box";
import { ThemeSwitch } from "@/components/custom/theme-switch";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { UserNav } from "@/components/user-nav";
import { useState } from "react";
import { Trans } from "react-i18next";
import { AllTab } from "./modules/tabs/all";
import { IdeaTab } from "./modules/tabs/idea";

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
        {/* start */}
        <div className='mb-2 flex items-center justify-between space-y-2'>
          <h1 className='text-2xl font-bold tracking-tight'>
            <Trans i18nKey={"greeting"}>Welcome!</Trans>
          </h1>
          <div className='flex items-center space-x-2'>
            <Button>Download</Button>
          </div>
        </div>
        {/* end */}
        <Tabs orientation='vertical' defaultValue='all' className='space-y-4'>
          <div>
            <TabsList>
              <TabsTrigger value='all'>All</TabsTrigger>
              <TabsTrigger value='idea'>Idea</TabsTrigger>
            </TabsList>
          </div>
          <AllTab />
          <IdeaTab />
        </Tabs>
      </Layout.Body>
    </Layout>
  );
};

export default Documents;
