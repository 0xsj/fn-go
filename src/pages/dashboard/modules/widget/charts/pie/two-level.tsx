import React, { PureComponent } from "react";
import { PieChart, Pie, ResponsiveContainer, Cell } from "recharts";

const data = [
  { name: "18-25", value: 150 },
  { name: "26-35", value: 250 },
  { name: "36-45", value: 200 },
  { name: "46-55", value: 180 },
  { name: "56+", value: 120 },
];

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#AF19FF"];

export default class EmployeeDemographicsChart extends PureComponent {
  render() {
    return (
      <ResponsiveContainer width='100%' height={400}>
        <PieChart>
          <Pie
            data={data}
            dataKey='value'
            cx='50%'
            cy='50%'
            outerRadius={80}
            fill='#8884d8'
            label
          >
            {data.map((entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
        </PieChart>
      </ResponsiveContainer>
    );
  }
}
