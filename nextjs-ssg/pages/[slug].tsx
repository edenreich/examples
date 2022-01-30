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

type Request = {
  slug: string | string[];
};

type Page = {
  slug: string;
  title: string;
  content: string;
};

const mockedPages: Page[] = [
  {
    slug: "about",
    title: "About",
    content: "About page",
  },
  {
    slug: "contact",
    title: "Contact",
    content: "Contact page",
  },
];

export const getStaticPaths: GetStaticPaths = async (): Promise<
  GetStaticPathsResult<Request>
> => {
  // Get all existing paths from an API for generating static pages ahead of time
  // Probably use Contentful API to get this data
  // But for now let's mock that data
  const allPossiblePaths = mockedPages.map((page) => ({
    params: { slug: page.slug },
  }));
  return {
    paths: allPossiblePaths,
    fallback: false,
  };
};

export const getStaticProps: GetStaticProps = async (
  context: GetStaticPropsContext
): Promise<GetStaticPropsResult<Page>> => {
  // Lookup for a page by slug, this is where you'd use Contentful API
  const page: Page = mockedPages.filter(
    (page) => page.slug === context.params?.slug
  )[0];

  return {
    props: {
      ...page,
    },
  };
};

const Blog: NextPage<Page> = ({ title, content }: PropsWithChildren<Page>) => {
  return (
    <div>
      <Head>
        <title>{title}</title>
        <meta name="description" content={title} />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>{content}</main>
      <footer></footer>
    </div>
  );
};

export default Blog;
