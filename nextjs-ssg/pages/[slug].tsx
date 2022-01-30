import type {
  NextPage,
  GetStaticProps,
  GetStaticPaths,
  GetStaticPropsContext,
  GetStaticPropsResult,
  GetStaticPathsResult,
} from "next";
import Head from "next/head";
import { PropsWithChildren } from "react";
import { promises } from "fs";

type Request = {
  slug: string | string[];
};

type Page = {
  slug: string;
  title: string;
  content: string;
};

const getPages = async (): Promise<Page[]> => {
  return JSON.parse(
    await promises.readFile("./data/pages.json", {
      encoding: "utf-8",
    })
  ) as unknown as Page[];
};

export const getStaticPaths: GetStaticPaths = async (): Promise<
  GetStaticPathsResult<Request>
> => {
  // Get all existing paths from an API for generating static pages ahead of time
  // Probably use Contentful API to get this data
  // But for now let's mock that data
  const mockedPages: Page[] = await getPages();
  const allPossiblePaths = mockedPages.map((page) => ({
    params: { slug: page.slug },
  }));
  return {
    paths: allPossiblePaths,
    fallback: true,
  };
};

export const getStaticProps: GetStaticProps = async (
  context: GetStaticPropsContext
): Promise<GetStaticPropsResult<Page>> => {
  // Lookup for a page by slug, this is where you'd use Contentful API
  const mockedPages: Page[] = await getPages();
  const page: Page = mockedPages.filter(
    (page) => page.slug === context.params?.slug
  )[0];
  console.log("called");
  return {
    props: {
      ...page,
    },
    revalidate: 30,
  };
};

const Blog: NextPage<Page> = (page: PropsWithChildren<Page>) => {
  if (!page) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <Head>
        <title>{page.title}</title>
        <meta name="description" content={page.title} />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>{page.content}</main>
      <footer></footer>
    </div>
  );
};

export default Blog;
