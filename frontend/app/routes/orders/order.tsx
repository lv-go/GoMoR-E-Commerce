import { DISPATCH_ACTION, PayPalButtons, usePayPalScriptReducer } from "@paypal/react-paypal-js";
import type { CreateOrderData, CreateOrderActions, OnApproveData, OnApproveActions } from "@paypal/paypal-js";
import { useEffect } from "react";
import { Link, useParams } from "react-router";
import { toast } from "react-hot-toast";
import Loader from "~/components/Loader";
import Messsage from "~/components/Message";
import { useFirebaseAuth } from "~/FirebaseAuthContext";
import {
  useGetById,
  useDeliverOrderMutation,
  // useGetPaypalClientIdQuery,
  usePayOrderMutation,
  newOrder
} from "~/hooks/orders";

export default function Order() {
  const { id: orderId } = useParams();

  const { data: order = newOrder(), refetch, isLoading, error } = useGetById(orderId || "");

  const { mutateAsync: payOrder, isPending: loadingPay } = usePayOrderMutation();
  const { mutateAsync: deliverOrder, isPending: loadingDeliver } = useDeliverOrderMutation();
  const { user } = useFirebaseAuth();

  const [{ isPending }] = usePayPalScriptReducer();

  async function onApprove(data: OnApproveData, actions: OnApproveActions) {
    const details = await actions.order?.capture()
    try {
      await payOrder({
        id: orderId || "", data: {
          paymentResult: {
            id: details?.id,
            status: details?.status,
            update_time: details?.update_time,
            email_address: details?.payer?.email_address,
          }
        }
      });
      refetch();
      toast.success("Order is paid");
    } catch (error: any) {
      toast.error(error?.data?.message || error.message);
    }
  }

  function createOrder(data: CreateOrderData, actions: CreateOrderActions) {
    return actions.order
      .create({
        intent: "CAPTURE",
        purchase_units: [{
          amount: {
            value: order.totalPrice.toFixed(2),
            currency_code: "USD"
          }
        }],
      })
      .then((orderID) => {
        console.log("Order ID: ", orderID);
        return orderID;
      });
  }

  function onError(err: any) {
    toast.error(err.message);
  }

  const deliverHandler = async () => {
    await deliverOrder(orderId || "");
    refetch();
  };

  return isLoading ? (
    <Loader />
  ) : error ? (
    <Messsage variant="error">{error.message}</Messsage>
  ) : (
    <div className="container mx-auto flex flex-col p-4">
      <div className="">
        <div className="border gray-300 mt-5 mb-5">
          {order.orderItems.length === 0 ? (
            <Messsage>Order is empty</Messsage>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="border-b-2">
                  <tr>
                    <th className="p-2">Image</th>
                    <th className="p-2">Product</th>
                    <th className="p-2 text-center">Quantity</th>
                    <th className="p-2">Unit Price</th>
                    <th className="p-2">Total</th>
                  </tr>
                </thead>

                <tbody>
                  {order.orderItems.map((item, index) => (
                    <tr key={index}>
                      <td className="p-2">
                        <img
                          src={item.image}
                          alt={item.name}
                          className="w-16 h-16 object-cover"
                        />
                      </td>

                      <td className="p-2">
                        <Link to={`/product/${item.productId}`}
                          className="link link-primary"
                        >
                          {item.name}
                        </Link>
                      </td>

                      <td className="p-2 text-center">{item.quantity}</td>
                      <td className="p-2 text-center">{item.price}</td>
                      <td className="p-2 text-center">
                        $ {(item.quantity * item.price).toFixed(2)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      <div className="w-full md:w-1/3 mx-auto">
        <div className="mt-5 border-gray-300 pb-4 mb-4">
          <h2 className="text-xl font-bold mb-2">Shipping</h2>
          <p className="mb-4 mt-4">
            <strong className="text-pink-500">Order:</strong> {order._id}
          </p>

          <p className="mb-4">
            <strong className="text-pink-500">Name:</strong>{" "}
            {user?.displayName}
          </p>

          <p className="mb-4">
            <strong className="text-pink-500">Email:</strong> {user?.email}
          </p>

          <p className="mb-4">
            <strong className="text-pink-500">Address:</strong>{" "}
            {order.shippingAddress.address}, {order.shippingAddress.city}{" "}
            {order.shippingAddress.postalCode}, {order.shippingAddress.country}
          </p>

          <p className="mb-4">
            <strong className="text-pink-500">Method:</strong>{" "}
            {order.paymentMethod}
          </p>

          {order.isPaid ? (
            <Messsage variant="success">Paid on {order.paidAt}</Messsage>
          ) : (
            <Messsage variant="error">Not paid</Messsage>
          )}
        </div>

        <h2 className="text-xl font-bold mb-2 mt-[3rem]">Order Summary</h2>
        <div className="flex justify-between mb-2">
          <span>Items</span>
          <span>$ {order.itemsPrice}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span>Shipping</span>
          <span>$ {order.shippingPrice}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span>Tax</span>
          <span>$ {order.taxPrice}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span>Total</span>
          <span>$ {order.totalPrice}</span>
        </div>

        {!order.isPaid && (
          <div>
            {loadingPay && <Loader />}{" "}
            {isPending ? (
              <Loader />
            ) : (
              <div>
                <div>
                  <PayPalButtons
                    createOrder={createOrder}
                    onApprove={onApprove}
                    onError={onError}
                  ></PayPalButtons>
                </div>
              </div>
            )}
          </div>
        )}

        {loadingDeliver && <Loader />}
        {user && user.role === "admin" && order.isPaid && !order.isDelivered && (
          <div>
            <button
              type="button"
              className="bg-pink-500 text-white w-full py-2"
              onClick={deliverHandler}
            >
              Mark As Delivered
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
