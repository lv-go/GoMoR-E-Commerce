import { useEffect } from "react";
import { useDispatch } from "react-redux";
import { toast } from "react-toastify";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Link } from "react-router";
import z from "zod";
import { useFirebaseAuth } from "~/FirebaseAuthContext";
import Loader from "../../components/Loader";
import { useProfileMutation } from "../../redux/api/usersApiSlice";
import { updateEmail, updatePassword, updateProfile } from "firebase/auth";

const ProfileSchema = z.object({
  displayName: z.string().min(3, "Username must be at least 3 characters long"),
  email: z.email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters long"),
  confirmPassword: z.string().min(6, "Confirm password must be at least 6 characters long"),
});

type Profile = z.infer<typeof ProfileSchema>;

export default function ProfilePage() {
  const { register, handleSubmit, formState: { errors }, reset } = useForm<Profile>({
    resolver: zodResolver(ProfileSchema),
  });

  const { user: userInfo } = useFirebaseAuth();

  useEffect(() => {
    if (userInfo) {
      reset({
        displayName: userInfo.displayName || "",
        email: userInfo.email || "",
      });
    }
  }, [userInfo, reset]);

  const submitHandler = async (data: Profile) => {
    const { displayName, email, password, confirmPassword } = data;
    if (password !== confirmPassword) {
      toast.error("Passwords do not match");
    } else {
      try {
        if (!userInfo) {
          toast.error("User not found");
          return;
        }
        updateProfile(userInfo, {
          displayName,
        })
        if (email !== userInfo.email) {
          await updateEmail(userInfo, email);
        }
        if (password && password !== "" && password !== confirmPassword) {
          await updatePassword(userInfo, password);
        }
        toast.success("Profile updated successfully");
      } catch (err: any) {
        toast.error(err?.data?.message || err.error);
      }
    }
  };

  return (
    <div className="container mx-auto p-4 _mt-[10rem]">
      <div className="flex justify-center align-center md:flex md:space-x-4">
        <div className="md:w-1/3">
          <h2 className="text-2xl font-semibold mb-4">Update Profile</h2>
          <form onSubmit={handleSubmit(submitHandler)}>
            <div className="mb-4">
              <label className="block text-white mb-2">Name</label>
              <input
                type="text"
                placeholder="Enter name"
                className="form-input p-4 rounded-sm w-full"
                {...register("displayName")}
              />
              {errors.displayName && <p className="text-red-500">{errors.displayName.message}</p>}
            </div>

            <div className="mb-4">
              <label className="block text-white mb-2">Email Address</label>
              <input
                type="email"
                placeholder="Enter email"
                className="form-input p-4 rounded-sm w-full"
                {...register("email")}
              />
              {errors.email && <p className="text-red-500">{errors.email.message}</p>}
            </div>

            <div className="mb-4">
              <label className="block text-white mb-2">Password</label>
              <input
                type="password"
                placeholder="Enter password"
                className="form-input p-4 rounded-sm w-full"
                {...register("password")}
              />
              {errors.password && <p className="text-red-500">{errors.password.message}</p>}
            </div>

            <div className="mb-4">
              <label className="block text-white mb-2">Confirm Password</label>
              <input
                type="password"
                placeholder="Confirm password"
                className="form-input p-4 rounded-sm w-full"
                {...register("confirmPassword")}
              />
              {errors.confirmPassword && <p className="text-red-500">{errors.confirmPassword.message}</p>}
            </div>

            <div className="flex justify-between">
              <button
                type="submit"
                className="bg-pink-500 text-white py-2 px-4 rounded hover:bg-pink-600"
              >
                Update
              </button>

              <Link
                to="/user-orders"
                className="bg-pink-600 text-white py-2 px-4 rounded hover:bg-pink-700"
              >
                My Orders
              </Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
