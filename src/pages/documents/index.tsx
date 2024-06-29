import { Layout } from "@/components/custom/layout"
import { useState } from "react"

interface Props {}

export const Documents: React.FC<Props> = () => {
  return (
    <Layout>
      <Layout.Header>
        <div>
          <h1>Logo</h1>
        </div>
      </Layout.Header>
      <Layout.Body>
        <div></div>
      </Layout.Body>
    </Layout>
  )
}

export default Documents
