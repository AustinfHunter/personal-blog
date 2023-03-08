import './App.scss';
import {
  createBrowserRouter,
  Outlet,
  RouterProvider,
} from "react-router-dom";
import Header from './components/Header';
import Footer from './components/Footer';
import Blog from './pages/Blog';
import Home from './pages/Home';
import Error from './pages/Error'
import Post from './pages/Post';
import Auth from './pages/Auth';
import SignIn from './pages/SignIn';
import PostEditor from './pages/PostEditor';
import PostManager from './pages/PostManager';

const router = createBrowserRouter([
  {
    path: "/",
    element: <DefaultLayout />,
    children: [
      {
        path: "/",
        element: <Home />
      },
      {
        path: "/blog",
        element: <Blog />,
      },
      {
        path: "/posts/:post",
        loader: async({params}) => {
          return params.post
        },
        element: <Post />
      }
    ],
    errorElement: <Error />,
  },
  {
    path: "/signin",
    element: <SignIn />
  },
  {
    path: "/admin",
    element: <Auth ><Admin/></Auth>,
    children: [
      {
        path: "/admin/editor",
        element: <PostEditor />
      },
      {
        path: "/admin/postmanager",
        element: <PostManager />
      }
    ]
  }
]);

function DefaultLayout(){
  return(
    <>
      <Header links={[{name: "Home", url: "/"}, {name:"Blog", url:"/blog"}]} logoContent={"Austin Hunter"}/>
      <Outlet className="outlet" />
      <Footer />
    </>
  )
}

function Admin() {
  return(
  <>
  <Header links={[{name: "Home", url: "/"}, {name:"Editor", url:"/admin/editor"}, {name:"Manage Posts", url:"/admin/postmanager"}]} logoContent={"Austin Hunter"}/>
  <Outlet className="outlet" />
  <Footer />
</>
  )
}

function App() {
  return (
    <div className="App">
      <RouterProvider router={router} />
    </div>
  );
}

export default App;
