import Table from "~/components/Table";
import { useCreate, useDeleteById, useGetPage, useUpdate } from "~/hooks/products";
import { getId, newProduct } from "~/schemas/product";
import { useGetPage as useCategoriesGetPage } from "~/hooks/categories";
import { useMemo } from "react";

export default function AllProducts() {
  const { data, isLoading, error, isFetched, isFetching } = useCategoriesGetPage({ offset: 0, limit: 100 });
  const categories = data?.items ?? [];

  const categoriesMap = useMemo(
    () => new Map(categories.map((category) =>
      [category._id, category])
    ),
    [categories]
  );

  return <Table
    header="Manage Products"
    newItem={newProduct}
    getItemId={getId}
    useGetPage={useGetPage}
    useCreate={useCreate}
    useUpdate={useUpdate}
    useDeleteById={useDeleteById}
    columns={[{
      key: "name",
      header: "Name",
      content: (item) => item.name,
      editor: (register) => <input className="input input-bordered" {...register} />,
    }, {
      key: "image",
      header: "Image",
      content: (item) => item.image,
      editor: (register) => <input className="input input-bordered" {...register} />,
    }, {
      key: "brand",
      header: "Brand",
      content: (item) => item.brand,
      editor: (register) => <input className="input input-bordered" {...register} />,
    }, {
      key: "category",
      header: "Category",
      content: (item) => categoriesMap.get(item.categoryId)?.name,
      editor: (register) => <select className="select select-bordered" {...register}>
        {categories.map((category) => (
          <option key={category._id} value={category._id}>{category.name}</option>
        ))}
      </select>,
    }, {
      key: "description",
      header: "Description",
      content: (item) => item.description,
      editor: (register) => <input className="input input-bordered" {...register} />,
    }, {
      key: "price",
      header: "Price",
      content: (item) => item.price,
      registerOptions: { min: 0, valueAsNumber: true },
      editor: (register) => <input type="number" className="input input-bordered" {...register} />,
    }, {
      key: "countInStock",
      header: "Count In Stock",
      content: (item) => item.countInStock,
      registerOptions: { min: 0, valueAsNumber: true },
      editor: (register) => <input type="number" className="input input-bordered" {...register} />,
    },]} />
}
