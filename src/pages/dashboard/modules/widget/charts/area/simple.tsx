import React from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  { week: "Week 1", ordersServed: 250, satisfactionScore: 85 },
  { week: "Week 2", ordersServed: 300, satisfactionScore: 88 },
  { week: "Week 3", ordersServed: 280, satisfactionScore: 86 },
  { week: "Week 4", ordersServed: 320, satisfactionScore: 90 },
  { week: "Week 5", ordersServed: 310, satisfactionScore: 87 },
];

export const EmployeeShiftPerformanceChart = () => {
  return (
    <div>
      <ResponsiveContainer width='100%' height={400}>
        <AreaChart
          data={data}
          margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='week' />
          <YAxis />
          <Tooltip />
          <Area
            type='monotone'
            dataKey='ordersServed'
            stackId='1'
            stroke='#8884d8'
            fill='#8884d8'
          />
          <Area
            type='monotone'
            dataKey='satisfactionScore'
            stackId='1'
            stroke='#82ca9d'
            fill='#82ca9d'
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};
