import Table from "~/components/Table";
import { useCreate, useDeleteById, useGetPage, useUpdate } from "~/hooks/categories";
import { getId, newCategory, type Category } from "~/schemas/category";

export default function CategoryList() {

  return (
    <Table
      header="Manage Categories"
      newItem={newCategory}
      getItemId={getId}
      useGetPage={useGetPage}
      useCreate={useCreate}
      useUpdate={useUpdate}
      useDeleteById={useDeleteById}
      columns={[
        {
          key: "name",
          header: "Name",
          content: (item) => item.name,
          editor: (register) => <input className="input input-bordered" {...register} />,
        },
      ]}
    />
  );
}
