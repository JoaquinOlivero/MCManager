import '../styles/globals.css'
import type { AppProps } from 'next/app'
import Head from 'next/head'
import { DataProvider } from '../contexts/DataContext'
import Layout from '../components/Layout/Layout'

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <div>
      <Head>
        <title>MCManager</title>
        <meta name="description" content="Manage your minecraft server" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <DataProvider>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </DataProvider>
    </div>

  )
}

export default MyApp
