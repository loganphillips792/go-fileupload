import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import "./index.css";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { imagesLoader } from "./FileList.jsx";

const router = createBrowserRouter([
    {
        path: "/",
        // element: <Root />,
        // loader: rootLoader,
        children: [
            {
                path: "app",
                element: <App />,
                loader: imagesLoader,
            },
        ],
    },
]);

ReactDOM.createRoot(document.getElementById("root")).render(<RouterProvider router={router} />);

// ReactDOM.createRoot(document.getElementById('root')).render(
//   <React.StrictMode>
//      {/* https://reactrouter.com/docs/en/v6/getting-started/tutorial */}
//      <BrowserRouter>
//       <Routes>
//         <Route path="/app" element={<App />}></Route>
//       </Routes>
//     </BrowserRouter>
//   </React.StrictMode>,
// )
